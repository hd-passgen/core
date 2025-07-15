# hd-pasgen

## Функционал

- Создание нового master-пароля
- Внесение существующего master-пароля
- Сохранение существующего master-пароля в TPM
- Сохраненение пароля для определенного сервиса
- Извлечение пароля для определенноо сервиса
- Получение пароля по идентификатору без необходимости вводить его (идентификатор) каждый раз
- Копирование извлеченного пароля сразу в буфер обмена без вывода в командную строку

## MVP

- Stateless командная утилита для получения пароля к необходимому сервису
- Дистрибьюция бинарного файла через apt менеджер пакетов

## Notes

### Alternatives

- [lesspass](https://github.com/lesspass/lesspass)
- [pass](https://www.passwordstore.org/)

### TPM

Software emulated TPM storage

```sh
sudo apt install swtpm swtpm-tools libtpms-dev tpm2-tools
```

```sh
# create dir to store tpm state
mkdir -p /tmp/tpmstate

# start swtpm socket server
swtpm socket \
  --tpmstate dir=/tmp/tpmstate \
  --ctrl type=unixio,path=/tmp/tpmstate/swtpm-sock \
  --tpm2 \
  --log level=20

# start character device
sudo swtpm chardev \
  --tpmstate dir=/tmp/tpmstate \
  --vtpm-proxy \
  --tpm2

# test
export TPM2TOOLS_TCTI="swtpm:path=/tmp/tpmstate/swtpm-sock"
tpm2_getrandom 10
tpm2_pcrread
```
