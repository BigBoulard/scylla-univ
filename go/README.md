# scylla-mutants 01

## Building the app image

```sh
docker build -t mutants .
[+] Building 0.6s (13/13) FINISHED
```

## Launch the app

```sh
docker run -d --net=scylla-univ_web --name mutants mutants
```

| option                | description                 |
| --------------------- | --------------------------- |
| -d                    | detached                    |
| --net=scylla-univ_web | use network scylla-univ_web |
| --name mutants        | name of the container       |
| mutants               | name of the image to use    |

## Connect to the app container

```sh
docker exec -it mutants sh
```
