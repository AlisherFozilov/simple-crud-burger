package burgers

const getIdNamePriceSQL = "SELECT id, name, price FROM burgers WHERE removed = FALSE"
const insertNamePriceSQL = `INSERT INTO burgers (name, price) VALUES ($1, $2);`
const setRemovedTrueByIdSQL = `UPDATE burgers SET removed = true WHERE id = $1;`