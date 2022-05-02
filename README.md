# img-receipt

이미지를 영수증프린터로 출력하는 프로그램

![example](_img/example.jpg)

## usage

```bash
GOARCH=arm GOOS=linux go build && \
scp img-receipt pi@rpi-pos.local:~/
```

```bash
curl -F "img=@./_img/Lenna.png" -k http://rpi-pos.local:8080/upload
```

## Reference
- [SEWOO, SLK-TS100 제품소개](https://www.miniprinter.com/ko/product/view.do?SEQ=159)