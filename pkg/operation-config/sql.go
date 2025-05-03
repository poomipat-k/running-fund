package operationConfig

const getOperationConfigSQL = `SELECT allow_new_project FROM operation_config ORDER BY id DESC LIMIT 1;`
