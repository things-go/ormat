ormat build \
   -i test_data1.sql \
   -i test_data2.sql \
   --out . \
   --hasColumn \
   --enableSea \
   --merge true \
   --package testdata \
   --options go_package=github.com/things-go/ormat/testdata
# --hasAssist
