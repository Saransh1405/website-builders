import { useState } from "react";
import { motion } from "framer-motion";
import { File, Folder, ChevronRight } from "lucide-react";

interface FileNode {
  name: string;
  type: "file" | "folder";
  children?: FileNode[];
  content?: string;
}

interface FileExplorerProps {
  files: FileNode[];
  onFileSelect: (file: FileNode) => void;
  selectedFile: FileNode | null;
}

function FileItem({
  node,
  depth,
  onFileSelect,
  selectedFile,
}: {
  node: FileNode;
  depth: number;
  onFileSelect: (file: FileNode) => void;
  selectedFile: FileNode | null;
}) {
  const isSelected = selectedFile?.name === node.name && node.type === "file";

  return (
    <div>
      <button
        onClick={() => node.type === "file" && onFileSelect(node)}
        className={`w-full flex items-center gap-2 px-3 py-2 text-sm transition-colors hover:bg-accent ${
          isSelected
            ? "bg-primary/10 text-primary"
            : "text-foreground"
        }`}
        style={{ paddingLeft: `${depth * 16 + 12}px` }}
      >
        {node.type === "folder" ? (
          <>
            <ChevronRight className="h-4 w-4 text-muted-foreground rotate-90" />
            <Folder className="h-4 w-4 text-primary" />
          </>
        ) : (
          <>
            <span className="w-4" />
            <File className="h-4 w-4 text-muted-foreground" />
          </>
        )}
        <span className="truncate font-mono text-xs">{node.name}</span>
      </button>

      {node.type === "folder" && node.children && (
        <div>
          {node.children.map((child, index) => (
            <FileItem
              key={`${child.name}-${index}`}
              node={child}
              depth={depth + 1}
              onFileSelect={onFileSelect}
              selectedFile={selectedFile}
            />
          ))}
        </div>
      )}
    </div>
  );
}

export function FileExplorer({
  files,
  onFileSelect,
  selectedFile,
}: FileExplorerProps) {
  return (
    <div className="h-full flex flex-col">
      <div className="p-4 border-b border-border">
        <h2 className="text-lg font-semibold text-foreground">Files</h2>
        <p className="text-sm text-muted-foreground mt-1">
          Click a file to view contents
        </p>
      </div>

      <div className="flex-1 overflow-y-auto py-2">
        {files.map((node, index) => (
          <motion.div
            key={`${node.name}-${index}`}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: index * 0.05 }}
          >
            <FileItem
              node={node}
              depth={0}
              onFileSelect={onFileSelect}
              selectedFile={selectedFile}
            />
          </motion.div>
        ))}
      </div>
    </div>
  );
}
