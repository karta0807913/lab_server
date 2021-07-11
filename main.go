package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/karta0807913/go_server_utils/serverutil"
	"github.com/karta0807913/lab_server/model"
	"github.com/karta0807913/lab_server/route"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	switch sql := WebsiteConfig.sql.(type) {
	case mysqlConfig:
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			sql.account,
			sql.password,
			sql.host,
			sql.port,
			sql.database,
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			return nil, err
		}
	case sqliteConfig:
		db, err = model.CreateSqliteDB(sql.filepath)
		if err != nil {
			return nil, err
		}
	}
	err = model.InitDB(db)
	if err != nil {
		return nil, err
	}
	return db, err
}

// if id not found, create a new calendar
func InitGoogleCalendar(calendarID string) (string, error) {
	dateRegexBegin := regexp.MustCompile(`^0*([0-9]{3})年0*([0-9]+)月0*([0-9]+)日`)
	dateRegexEnd := regexp.MustCompile(`0*([0-9]{3})年0*([0-9]+)月0*([0-9]+)日$`)
	ctx := context.Background()
	svc, err := calendar.NewService(ctx, option.WithCredentialsFile(WebsiteConfig.googleAuthFile))
	if err != nil {
		return "", err
	}

	var targetCalendar *calendar.Calendar
	targetCalendar, err = svc.Calendars.Get(calendarID).Do()
	if err != nil {
		targetCalendar, err = svc.Calendars.Insert(&calendar.Calendar{
			Summary: "科技部計畫時程",
		}).Do()
	}
	_, err = svc.Acl.Insert(targetCalendar.Id, &calendar.AclRule{
		Role: "reader",
		Scope: &calendar.AclRuleScope{
			Type: "default",
		},
	}).Do()
	if err != nil {
		return "", err
	}
	eventsMap := make(map[string]*calendar.Event)
	events, err := svc.Events.List(targetCalendar.Id).Do()
	log.Println(len(events.Items))
	for _, event := range events.Items {
		eventsMap[event.Summary] = event
	}
	res, err := http.Get("https://www.most.gov.tw/folksonomy/rfpList?pageSize=400&l=ch")
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}
	doc.Find(".plan_table tbody tr").Each(func(i int, s *goquery.Selection) {
		// get event date
		date := strings.ReplaceAll(strings.ReplaceAll(s.Find("td[headers='activityStartDate']").Text(), " ", ""), "\n", "")
		beginRegexDate := dateRegexBegin.FindStringSubmatch(date)
		if len(beginRegexDate) != 4 {
			log.Println("cannot get event date")
			return
		}
		beginYear, err := strconv.Atoi(beginRegexDate[1])
		if err != nil {
			log.Println("fetch date error", beginRegexDate)
			return
		}
		endRegexDate := dateRegexEnd.FindStringSubmatch(date)
		if len(endRegexDate) != 4 {
			endRegexDate = beginRegexDate
		}
		endYear, err := strconv.Atoi(endRegexDate[1])
		if err != nil {
			endRegexDate = beginRegexDate
			endYear = beginYear
		}

		paddingDate := func(s string) string {
			for i := len(s); i < 2; i++ {
				s = "0" + s
			}
			return s
		}

		endYear += 1911
		beginYear += 1911
		beginDate := strconv.Itoa(beginYear) + "-" + paddingDate(beginRegexDate[2]) + "-" + paddingDate(beginRegexDate[3])
		endDate := strconv.Itoa(endYear) + "-" + paddingDate(endRegexDate[2]) + "-" + paddingDate(endRegexDate[3])

		summary, ok := s.Find(".link_icon").Attr("title")
		if !ok {
			log.Println("cannot get event title")
			return
		}
		url, ok := s.Find(".link_icon").Attr("data-target-url")
		if !ok {
			log.Println("cannot get event url")
			url = "/folksonomy/rfpList"
		}

		// check if event exists
		if event, ok := eventsMap[summary]; ok {
			// update
			log.Println("update event", summary)
			event.Description = "https://www.most.gov.tw" + url
			event.Start = &calendar.EventDateTime{
				Date: beginDate,
			}
			event.End = &calendar.EventDateTime{
				Date: endDate,
			}
			svc.Events.Update(targetCalendar.Id, event.Id, event).Do()
		} else {
			log.Println("add event", summary)
			_, err = svc.Events.Insert(targetCalendar.Id, &calendar.Event{
				Summary:     summary,
				Description: "https://www.most.gov.tw" + url,
				Start: &calendar.EventDateTime{
					Date: beginDate,
				},
				End: &calendar.EventDateTime{
					Date: endDate,
				},
			}).Do()
		}
		if err != nil {
			log.Println(err)
		}
	})
	return targetCalendar.Id, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// calendar updater
	go func() {
		ID, err := InitGoogleCalendar(WebsiteConfig.calendarID)
		if err != nil {
			log.Println("update calendar failed, error is", err)
		} else {
			log.Println("calendar", ID, "updated")
		}
		// update calendar every 48 hours
		for {
			timer := time.NewTimer(48 * time.Hour)
			select {
			case <-timer.C:
				ID, err := InitGoogleCalendar(WebsiteConfig.calendarID)
				if err != nil {
					log.Println("update calendar failed, error is", err)
				} else {
					log.Println("calendar", ID, "updated")
				}
			}
		}
	}()

	db, err := InitDB()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := serverutil.NewGormStorage(db)
	if err != nil {
		log.Fatal(err)
	}

	server, err :=
		serverutil.NewGinServer(serverutil.ServerSettings{
			PrivateKeyPath: WebsiteConfig.privateKeyPath,
			ServerAddress:  WebsiteConfig.serverAddr,
			Db:             db,
			Storage:        storage,
			SessionName:    "session",
		})
	if err != nil {
		log.Fatal(err)
	}

	route.Route(route.RouteConfig{
		DB:         db,
		Server:     server,
		UploadPath: WebsiteConfig.uploadPath,
	})

	log.Printf("server listening on %s", WebsiteConfig.serverAddr)
	log.Fatal(server.Run(WebsiteConfig.serverAddr))
}
