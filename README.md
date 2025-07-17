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

TODO: сделать анализ функционала, выделить плюсы и минусы каждого решения

- [lesspass](https://github.com/lesspass/lesspass)
- [pass](https://www.passwordstore.org/)

### TPM

My machine (TP X13 gen 1 with AMD Ryzen 7 4750u) has only TPM 2.0 devices located on `/dev/tpm0` and `/dev/tpmrm0` (linux)

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

To use Secure Enclave on Mac, your app must have a registered App ID (`com.apple.application-identifier` entitlement). For more information, see [this thread](https://developer.apple.com/forums/thread/728150).

### Browser extensions for communicating with host-OS processes

- [https://developer.chrome.com/docs/extensions/develop/concepts/native-messaging]
- [https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Native_messaging]

### Password generator

Для генерации пароля необходимо иметь:
  - master пароль; 
  - mnemonic words [BIP39](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) (последовательность слов, разделённых пробелами - пока непонятно откуда их брать).

TODO: 
- [ ] Придумать, как сделать генерацию child пароля максимально устойчивой к коллизиям
- [ ] Мастер пароль нужно сохранять заранее, чтобы его не приходилось вводить
- [ ] Как поступаем с mnemonic words