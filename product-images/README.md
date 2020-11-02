# USAGE

## Upload

```
http POST :9090/images/1/frappuccino.png @frappuccino.png
```
Files will get stored in `BASE_PATH` (defaults to `./imagestore`).

## Download

```
http :9090/images/1/frappuccino.png
```
