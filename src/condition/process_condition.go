package condition


type ProcessCondition interface {
    Run(pid int, include_children bool) (float64 , error)
    Check(value float64, include_children bool) (bool, error)
    FormatValue(value float64) string
}