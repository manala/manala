package components

type Table struct {
	Rows []*TableRow
}

func (table *Table) AddRow(primary string, secondary string) {
	table.Rows = append(table.Rows, &TableRow{
		Primary:   primary,
		Secondary: secondary,
	})
}

type TableRow struct {
	Primary   string
	Secondary string
}
