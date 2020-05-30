package utils

/***
 * web server 永远都不应该crash
 * 启动新的routine，如果没有recover保证，一旦遇到panic，影响会扩散出去
 * crash的就是整个应用
 */
func CallFuncInNewRecoveryRoutine(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ErrorLogger.Logf(string(Stack(3)))
			}
		}()
		f()
	}()
}
