package operationConfig

const getOperationConfigSQL = `SELECT allow_new_project FROM operation_config ORDER BY id DESC LIMIT 1;`

const updateOperationConfigSQL = `UPDATE operation_config SET allow_new_project = $1 WHERE id = (select id FROM operation_config ORDER BY id DESC limit 1);`
