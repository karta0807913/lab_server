if [ -z "$DEBUG" ]
then
    DEBUG="-v --trace-ascii  /dev/stdout"
fi

curl $HOST/file/upload $OPT --form json='{ "filename": "test.jpg" };type=application/json' --form upload='@upload.bash' $DEBUG
