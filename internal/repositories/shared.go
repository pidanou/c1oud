package repositories

func isValidOrderBy(orderBy string) bool {
	allowedColumns := map[string]bool{
		"account_id":    true,
		"resource_name": true,
	}
	return allowedColumns[orderBy]
}
