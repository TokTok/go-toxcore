[![Build Status](https://travis-ci.org/kitech/go-toxcore.svg?branch=master)](https://travis-ci.org/kitech/go-toxcore)
[![GoDoc](https://godoc.org/github.com/kitech/go-toxcore?status.svg)](https://godoc.org/github.com/kitech/go-toxcore)

## go-toxcore
The golang bindings for libtoxcore 


### Installation

    # fetch libtoxcore if necessary
    # see https://github.com/irungentoo/toxcore/blob/master/INSTALL.md
    go get github.com/TokTok/go-toxcore-c


### Examples

    import "github.com/kitech/go-toxcore"

    // use custom options
    opt := tox.NewToxOptions()
    t := tox.NewTox(opt)
    av := tox.NewToxAv(t)
    
    // use default options
    t := tox.NewTox(nil)
    av := tox.NewToxAv(t)

### Tests

    go test -v -covermode count
    

Contributing
------------
1. Fork it
2. Create your feature branch (``git checkout -b my-new-feature``)
3. Commit your changes (``git commit -am 'Add some feature'``)
4. Push to the branch (``git push origin my-new-feature``)
5. Create new Pull Request

