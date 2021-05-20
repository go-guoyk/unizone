# sugar

abstracted interface for zap.SugaredLogger

## Usage

For library developer, you could simply import `go.guoyk.net/sugar` and ask user to pass in a`sugar.Logger`
implementation as argument.

For user, you should use `sugar_zap.Wrap` to wrap a existed `zap.Logger`.

## Credits

Guo Y.K., MIT License
