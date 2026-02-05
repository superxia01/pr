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
  cell?: (props: { row: T; value: any; isEditing?: boolean; onStartEdit?: () => void; onUpdate?: (value: any) => void }) => React.ReactNode
  sortable?: boolean
  filterable?: boolean
  editable?: boolean
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
  onCellUpdate?: (rowId: string, field: string, value: any) => void
  emptyMessage?: string
  className?: string
  // 服务端分页支持
  serverSide?: boolean
  total?: number
  onPageChange?: (pageIndex: number, pageSize: number) => void
  loading?: boolean
}

export function DataTable<T>({
  data,
  columns,
  searchable = true,
  pageSize = 10,
  pageSizeOptions = [10, 20, 30, 50],
  onRowClick,
  onCellUpdate,
  emptyMessage = "暂无数据",
  className,
  serverSide = false,
  total,
  onPageChange,
  loading = false,
}: DataTableProps<T>) {
  const [search, setSearch] = useState("")
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize,
  })
  const [sorting, setSorting] = useState<SortingState>(null)
  const [filterColumn, setFilterColumn] = useState<string | null>(null)
  const [filterValue, setFilterValue] = useState("")

  // 编辑状态
  const [editingCell, setEditingCell] = useState<{ rowId: string; field: string } | null>(null)

  // 处理分页变化
  const handlePaginationChange = (newPagination: PaginationState) => {
    setPagination(newPagination)
    if (serverSide && onPageChange) {
      onPageChange(newPagination.pageIndex, newPagination.pageSize)
    }
  }

  // 客户端模式：计算过滤后的数据
  const filteredData = useMemo(() => {
    if (serverSide) return data // 服务端模式不过滤

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
  }, [data, search, filterColumn, filterValue, columns, serverSide])

  // 客户端模式：排序
  const sortedData = useMemo(() => {
    if (serverSide) return data // 服务端模式不排序

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
  }, [filteredData, sorting, columns, serverSide, data])

  // 客户端模式：分页
  const paginatedData = useMemo(() => {
    if (serverSide) return data // 服务端模式直接使用传入的数据

    const start = pagination.pageIndex * pagination.pageSize
    return sortedData.slice(start, start + pagination.pageSize)
  }, [sortedData, pagination, serverSide, data])

  // 计算总页数
  const totalPages = Math.ceil((serverSide ? total || 0 : sortedData.length) / pagination.pageSize)

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

  // 获取行ID
  const getRowId = (row: T) => {
    return (row as any).id || JSON.stringify(row)
  }

  // 处理单元格更新
  const handleCellUpdate = (row: T, column: ColumnDef<T>, value: any) => {
    const rowId = getRowId(row)
    if (onCellUpdate && column.accessorKey) {
      const field = typeof column.accessorKey === "function"
        ? column.id
        : String(column.accessorKey)
      onCellUpdate(rowId, field, value)
      setEditingCell(null)
    }
  }

  // 检查是否正在编辑
  const isEditing = (row: T, column: ColumnDef<T>) => {
    const rowId = getRowId(row)
    const field = typeof column.accessorKey === "function"
      ? column.id
      : String(column.accessorKey)
    return editingCell?.rowId === rowId && editingCell?.field === field
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
      {/* 搜索和过滤 - 仅客户端模式 */}
      {!serverSide && searchable && (
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
            {loading ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  <div className="flex items-center justify-center">
                    <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-gray-900"></div>
                  </div>
                </TableCell>
              </TableRow>
            ) : paginatedData.length === 0 ? (
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
                  {columns.map((column) => {
                    const cellIsEditing = isEditing(row, column)
                    const cellValue = getCellValue(row, column)

                    return (
                      <TableCell key={column.id}>
                        {column.editable && !column.cell ? (
                          // 可编辑单元格（默认Input）
                          cellIsEditing ? (
                            <div className="flex items-center gap-2">
                              <Input
                                defaultValue={cellValue}
                                autoFocus
                                onBlur={(e) => handleCellUpdate(row, column, e.target.value)}
                                onKeyDown={(e) => {
                                  if (e.key === 'Enter') {
                                    handleCellUpdate(row, column, e.currentTarget.value)
                                  } else if (e.key === 'Escape') {
                                    setEditingCell(null)
                                  }
                                }}
                                className="h-8"
                              />
                            </div>
                          ) : (
                            <div
                              className="min-h-[2rem] cursor-pointer hover:bg-muted/50 rounded px-2 py-1 -mx-2"
                              onClick={() => setEditingCell({ rowId: getRowId(row), field: String(column.accessorKey) })}
                            >
                              {String(cellValue ?? "")}
                            </div>
                          )
                        ) : column.cell ? (
                          // 自定义渲染函数
                          column.cell({
                            row,
                            value: cellValue,
                            isEditing: cellIsEditing,
                            onStartEdit: () => setEditingCell({ rowId: getRowId(row), field: String(column.accessorKey) }),
                            onUpdate: (value) => handleCellUpdate(row, column, value),
                          })
                        ) : (
                          // 普通单元格
                          String(cellValue ?? "")
                        )}
                      </TableCell>
                    )
                  })}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {/* 分页 */}
      <div className="flex items-center justify-between px-2">
        <div className="text-sm text-muted-foreground">
          共 {serverSide ? total || 0 : sortedData.length} 条记录
        </div>

        <div className="flex items-center space-x-6 lg:space-x-8">
          <div className="flex items-center space-x-2">
            <p className="text-sm font-medium">每页</p>
            <select
              value={pagination.pageSize}
              onChange={(e) =>
                handlePaginationChange({
                  ...pagination,
                  pageSize: Number(e.target.value),
                  pageIndex: 0,
                })
              }
              className="h-8 w-[70px] rounded-md border border-input bg-background px-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              disabled={loading}
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
                handlePaginationChange({
                  ...pagination,
                  pageIndex: Math.max(0, pagination.pageIndex - 1),
                })
              }
              disabled={pagination.pageIndex === 0 || loading}
            >
              上一页
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() =>
                handlePaginationChange({
                  ...pagination,
                  pageIndex: Math.min(totalPages - 1, pagination.pageIndex + 1),
                })
              }
              disabled={pagination.pageIndex >= totalPages - 1 || loading}
            >
              下一页
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
