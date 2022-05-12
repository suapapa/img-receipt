# img-receipt

이미지를 영수증프린터로 출력하는 프로그램

![example](_img/example.jpg)

## usage

```bash
GOARCH=arm GOOS=linux go build && \
scp img-receipt pi@rpi-pos.local:~/
```

```bash
curl -F "img=@./_img/Lenna.png" http://rpi-pos.local:8080/upload
```

```bash
curl -X POST -d {"content": "hello world"} http://rpi-pos.local:8080/qr
```

## Reference

- [SEWOO, SLK-TS100 제품소개](https://www.miniprinter.com/ko/product/view.do?SEQ=159)
- [Barcode Contents](https://github.com/zxing/zxing/wiki/Barcode-Contents)
