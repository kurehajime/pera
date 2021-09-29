# pera

[![CircleCI](https://circleci.com/gh/kurehajime/pera.svg?style=svg)](https://circleci.com/gh/kurehajime/pera)

pera is the command to slide display the text in CLI.

![screenshot1](https://cloud.githubusercontent.com/assets/4569916/16450500/3ed1e1be-3e38-11e6-860a-084bf6d82b0f.gif)

---

## Installation

Download [here](https://github.com/kurehajime/pera/releases).

or 

Build.

```
go install github.com/kurehajime/pera@latest
```

---

## Usage

### [1] make text file.

```

first page.

---

second page.

---

third page.

```

 It will be a new page in three hyphens(---) on beginning of line.

---

### [2] do pera 

```
$ pera example.txt
```

Slide begins.

---

## More usage

If you run with no parameters, displays help slide.

```
$ pera
```

---

## License

[MIT](https://github.com/kurehajime/pera/blob/master/LICENSE)

---

![chaofan](https://cloud.githubusercontent.com/assets/4569916/16450501/3edc6b02-3e38-11e6-93c7-9cbd2a6c40f2.gif)
