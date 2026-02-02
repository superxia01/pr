import { useState, useMemo } from "react"
import {
  ChevronDown,
  ChevronUp,
  ChevronsUpDown,
  Search,
} from "lucide-react"
import { Button } from "./button"
import { Input } from "./input"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "./table"
import { cn } from "../../lib/utils"

// 类型定义
export type ColumnDef<T> = {
  id: string
  header: string
  accessorKey?: keyof T | ((row: T) => any)
  cell?: (props: { row: T; value: any }) => React.ReactNode
  sortable?: boolean
  filterable?: boolean
}

export type PaginationState = {
  pageIndex: number
  pageSize: number
}

export type SortingState = {
  id: string
  desc: boolean
} | null

type DataTableProps<T> = {
  data: T[]
  columns: ColumnDef<T>[]
  searchable?: boolean
  pageSize?: number
  pageSizeOptions?: number[]
  onRowClick?: (row: T) => void
  emptyMessage?: string
  className?: string
}

export function DataTable<T>({
  data,
  columns,
  searchable = true,
  pageSize = 10,
  pageSizeOptions = [10, 20, 30, 50],
  onRowClick,
  emptyMessage = "暂无数据",
  className,
}: DataTableProps<T>) {
  const [search, setSearch] = useState("")
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize,
  })
  const [sorting, setSorting] = useState<SortingState>(null)
  const [filterColumn, setFilterColumn] = useState<string | null>(null)
  const [filterValue, setFilterValue] = useState("")

  // 计算过滤后的数据
  const filteredData = useMemo(() => {
    return data.filter((item) => {
      // 全局搜索
      if (search) {
        const searchLower = search.toLowerCase()
        const matchSearch = columns.some((column) => {
          const value = column.accessorKey
            ? item[column.accessorKey as keyof T]
            : null
          return (
            value &&
            String(value)
              .toLowerCase()
              .includes(searchLower)
          )
        })
        if (!matchSearch) return false
      }

      // 列过滤
      if (filterColumn && filterValue) {
        const column = columns.find((col) => col.id === filterColumn)
        if (column && column.accessorKey) {
          const value = item[column.accessorKey as keyof T]
          if (!String(value).toLowerCase().includes(filterValue.toLowerCase())) {
            return false
          }
        }
      }

      return true
    })
  }, [data, search, filterColumn, filterValue, columns])

  // 排序
  const sortedData = useMemo(() => {
    if (!sorting) return filteredData

    return [...filteredData].sort((a, b) => {
      const column = columns.find((col) => col.id === sorting.id)
      if (!column || !column.accessorKey) return 0

      const aValue = column.accessorKey
        ? a[column.accessorKey as keyof T]
        : null
      const bValue = column.accessorKey
        ? b[column.accessorKey as keyof T]
        : null

      if (aValue === bValue) return 0
      if (aValue == null) return 1
      if (bValue == null) return -1

      const comparison = (aValue ?? 0) > (bValue ?? 0) ? 1 : -1
      return sorting.desc ? -comparison : comparison
    })
  }, [filteredData, sorting, columns])

  // 分页
  const paginatedData = useMemo(() => {
    const start = pagination.pageIndex * pagination.pageSize
    return sortedData.slice(start, start + pagination.pageSize)
  }, [sortedData, pagination])

  const totalPages = Math.ceil(sortedData.length / pagination.pageSize)

  // 获取单元格值
  const getCellValue = (row: T, column: ColumnDef<T>) => {
    if (column.accessorKey) {
      const value =
        typeof column.accessorKey === "function"
          ? column.accessorKey(row)
          : row[column.accessorKey]
      return value
    }
    return null
  }

  // 处理排序
  const handleSort = (columnId: string) => {
    const column = columns.find((col) => col.id === columnId)
    if (!column || !column.sortable) return

    setSorting((current) => {
      if (current?.id === columnId) {
        if (current.desc) {
          return null // 取消排序
        }
        return { id: columnId, desc: true } // 切换为降序
      }
      return { id: columnId, desc: false } // 升序
    })
  }

  return (
    <div className={cn("space-y-4", className)}>
      {/* 搜索和过滤 */}
      {searchable && (
        <div className="flex items-center gap-4">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              placeholder="搜索..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9"
            />
          </div>

          {/* 列过滤 */}
          {columns.filter((col) => col.filterable).length > 0 && (
            <div className="flex items-center gap-2">
              <select
                value={filterColumn || ""}
                onChange={(e) => {
                  setFilterColumn(e.target.value || null)
                  setFilterValue("")
                }}
                className="h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              >
                <option value="">全部列</option>
                {columns
                  .filter((col) => col.filterable)
                  .map((col) => (
                    <option key={col.id} value={col.id}>
                      {col.header}
                    </option>
                  ))}
              </select>

              {filterColumn && (
                <Input
                  placeholder="过滤值..."
                  value={filterValue}
                  onChange={(e) => setFilterValue(e.target.value)}
                  className="w-40"
                />
              )}
            </div>
          )}
        </div>
      )}

      {/* 表格 */}
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              {columns.map((column) => (
                <TableHead key={column.id}>
                  {column.sortable ? (
                    <Button
                      variant="ghost"
                      size="sm"
                      className="-ml-3 h-8 data-[state=open]:bg-accent"
                      onClick={() => handleSort(column.id)}
                    >
                      {column.header}
                      {sorting?.id === column.id ? (
                        sorting.desc ? (
                          <ChevronDown className="ml-2 h-4 w-4" />
                        ) : (
                          <ChevronUp className="ml-2 h-4 w-4" />
                        )
                      ) : (
                        <ChevronsUpDown className="ml-2 h-4 w-4" />
                      )}
                    </Button>
                  ) : (
                    <span>{column.header}</span>
                  )}
                </TableHead>
              ))}
            </TableRow>
          </TableHeader>
          <TableBody>
            {paginatedData.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  {emptyMessage}
                </TableCell>
              </TableRow>
            ) : (
              paginatedData.map((row, rowIndex) => (
                <TableRow
                  key={rowIndex}
                  className={cn(onRowClick && "cursor-pointer")}
                  onClick={() => onRowClick?.(row)}
                >
                  {columns.map((column) => (
                    <TableCell key={column.id}>
                      {column.cell
                        ? column.cell({
                            row,
                            value: getCellValue(row, column),
                          })
                        : String(getCellValue(row, column) ?? "")}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {/* 分页 */}
      <div className="flex items-center justify-between px-2">
        <div className="text-sm text-muted-foreground">
          共 {sortedData.length} 条记录
        </div>

        <div className="flex items-center space-x-6 lg:space-x-8">
          <div className="flex items-center space-x-2">
            <p className="text-sm font-medium">每页</p>
            <select
              value={pagination.pageSize}
              onChange={(e) =>
                setPagination({
                  ...pagination,
                  pageSize: Number(e.target.value),
                  pageIndex: 0,
                })
              }
              className="h-8 w-[70px] rounded-md border border-input bg-background px-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            >
              {pageSizeOptions.map((size) => (
                <option key={size} value={size}>
                  {size}
                </option>
              ))}
            </select>
          </div>

          <div className="flex w-[100px] items-center justify-center text-sm font-medium">
            第 {pagination.pageIndex + 1} / {totalPages || 1} 页
          </div>

          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() =>
                setPagination({
                  ...pagination,
                  pageIndex: Math.max(0, pagination.pageIndex - 1),
                })
              }
              disabled={pagination.pageIndex === 0}
            >
              上一页
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() =>
                setPagination({
                  ...pagination,
                  pageIndex: Math.min(totalPages - 1, pagination.pageIndex + 1),
                })
              }
              disabled={pagination.pageIndex >= totalPages - 1}
            >
              下一页
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
