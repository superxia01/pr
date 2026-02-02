import { useEditor, EditorContent } from "@tiptap/react"
import StarterKit from "@tiptap/starter-kit"
import Placeholder from "@tiptap/extension-placeholder"
import Link from "@tiptap/extension-link"
import {
  Bold,
  Italic,
  Underline,
  List,
  ListOrdered,
  Heading1,
  Heading2,
  Heading3,
  Link as LinkIcon,
  Undo,
  Redo,
} from "lucide-react"
import { Button } from "./button"
import { cn } from "../../lib/utils"
import { forwardRef } from "react"

interface RichTextEditorProps {
  value?: string
  onChange?: (value: string) => void
  placeholder?: string
  className?: string
  editable?: boolean
}

const MenuBar = ({ editor }: { editor: any }) => {
  if (!editor) {
    return null
  }

  return (
    <div className="flex flex-wrap items-center gap-1 rounded-t-md border border-b-0 border-input bg-muted/50 p-2">
      {/* 撤销/重做 */}
      <Button
        type="button"
        variant="ghost"
        size="sm"
        onClick={() => editor.chain().focus().undo().run()}
        disabled={!editor.can().undo()}
        title="撤销 (Ctrl+Z)"
      >
        <Undo className="h-4 w-4" />
      </Button>
      <Button
        type="button"
        variant="ghost"
        size="sm"
        onClick={() => editor.chain().focus().redo().run()}
        disabled={!editor.can().redo()}
        title="重做 (Ctrl+Y)"
      >
        <Redo className="h-4 w-4" />
      </Button>

      <div className="w-px bg-border mx-1" />

      {/* 标题 */}
      <Button
        type="button"
        variant={editor.isActive("heading", { level: 1 }) ? "secondary" : "ghost"}
        size="sm"
        onClick={() => editor.chain().focus().toggleHeading({ level: 1 }).run()}
        title="标题1"
      >
        <Heading1 className="h-4 w-4" />
      </Button>
      <Button
        type="button"
        variant={editor.isActive("heading", { level: 2 }) ? "secondary" : "ghost"}
        size="sm"
        onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
        title="标题2"
      >
        <Heading2 className="h-4 w-4" />
      </Button>
      <Button
        type="button"
        variant={editor.isActive("heading", { level: 3 }) ? "secondary" : "ghost"}
        size="sm"
        onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()}
        title="标题3"
      >
        <Heading3 className="h-4 w-4" />
      </Button>

      <div className="w-px bg-border mx-1" />

      {/* 格式化 */}
      <Button
        type="button"
        variant={editor.isActive("bold") ? "secondary" : "ghost"}
        size="sm"
        onClick={() => editor.chain().focus().toggleBold().run()}
        title="粗体 (Ctrl+B)"
      >
        <Bold className="h-4 w-4" />
      </Button>
      <Button
        type="button"
        variant={editor.isActive("italic") ? "secondary" : "ghost"}
        size="sm"
        onClick={() => editor.chain().focus().toggleItalic().run()}
        title="斜体 (Ctrl+I)"
      >
        <Italic className="h-4 w-4" />
      </Button>
      <Button
        type="button"
        variant={editor.isActive("underline") ? "secondary" : "ghost"}
        size="sm"
        onClick={() => editor.chain().focus().toggleUnderline().run()}
        title="下划线 (Ctrl+U)"
      >
        <Underline className="h-4 w-4" />
      </Button>

      <div className="w-px bg-border mx-1" />

      {/* 列表 */}
      <Button
        type="button"
        variant={editor.isActive("bulletList") ? "secondary" : "ghost"}
        size="sm"
        onClick={() => editor.chain().focus().toggleBulletList().run()}
        title="无序列表"
      >
        <List className="h-4 w-4" />
      </Button>
      <Button
        type="button"
        variant={editor.isActive("orderedList") ? "secondary" : "ghost"}
        size="sm"
        onClick={() => editor.chain().focus().toggleOrderedList().run()}
        title="有序列表"
      >
        <ListOrdered className="h-4 w-4" />
      </Button>

      <div className="w-px bg-border mx-1" />

      {/* 链接 */}
      <Button
        type="button"
        variant={editor.isActive("link") ? "secondary" : "ghost"}
        size="sm"
        onClick={() => {
          const url = window.prompt("输入链接地址:")
          if (url) {
            editor.chain().focus().setLink({ href: url }).run()
          }
        }}
        title="插入链接"
      >
        <LinkIcon className="h-4 w-4" />
      </Button>
      {editor.isActive("link") && (
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={() => editor.chain().focus().unsetLink().run()}
          title="移除链接"
        >
          取消链接
        </Button>
      )}
    </div>
  )
}

export const RichTextEditor = forwardRef<
  HTMLDivElement,
  RichTextEditorProps
>(({ value, onChange, placeholder, className, editable = true }, ref) => {
  const editor = useEditor({
    extensions: [
      StarterKit.configure({
        heading: {
          levels: [1, 2, 3],
        },
      }),
      Placeholder.configure({
        placeholder: placeholder || "请输入内容...",
      }),
      Link.configure({
        openOnClick: false,
        HTMLAttributes: {
          class: "text-primary underline",
        },
      }),
    ],
    content: value || "",
    editable,
    onUpdate: ({ editor }) => {
      onChange?.(editor.getHTML())
    },
    editorProps: {
      attributes: {
        class:
          "prose prose-sm sm:prose lg:prose-lg xl:prose-xl focus:outline-none max-w-none min-h-[200px] p-4",
      },
    },
  })

  return (
    <div ref={ref} className={cn("rounded-md border", className)}>
      {editable && <MenuBar editor={editor} />}
      <EditorContent editor={editor} />
    </div>
  )
})

RichTextEditor.displayName = "RichTextEditor"
