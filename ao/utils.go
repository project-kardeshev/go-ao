package ao

func Coalesce[T any](value interface{}, defaultValue any) T {
    if value != nil {
        return value.(T)
    }
    return defaultValue.(T)
}