package enum

const (
	OPERATE_SUCCESS int = 0 // 操作成功
	OPERATE_FAILED  int = 1 // 操作失败
)

var CodeMap map[int]string = map[int]string{
	OPERATE_SUCCESS: "操作成功",
	OPERATE_FAILED:  "操作失败",
}
