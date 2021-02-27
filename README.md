# Резервное копирование с litres.

Я купил 500+ книг с litres.ru. Внезапно мне показалось, что неплохо было-бы сделать резервную копию всех моих накоплений.

Все книги выкачиваются параллельно...

**Это не "бесплатная качалка платных книг". Скачивает только те книги, что есть в разделе "мои книги" на litres.ru**

_Базируется на litres api версии 3.31_

## Dependencies
- github.com/cheggaaa/pb

## Как пользоваться

```
  -b    Show progress bar
  -d    print lots of debugging information
  -f string
        Downloading format. 'list' for available (default "list")
  -l string
        The directory where the books will be saved (default "/tmp")
  -p string
        password (default "your_password")
  -u string
        username (default "your_login")
  -v    be verbose (this is the default)
```