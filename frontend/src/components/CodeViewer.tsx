import { motion } from "framer-motion";

interface CodeViewerProps {
  fileName: string;
  content: string;
}

export function CodeViewer({ fileName, content }: CodeViewerProps) {
  const lines = content.split("\n");

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="h-full flex flex-col"
    >
      <div className="px-4 py-3 border-b border-border bg-card">
        <span className="font-mono text-sm text-foreground">{fileName}</span>
      </div>

      <div className="flex-1 overflow-auto bg-code-bg">
        <pre className="p-4 text-sm">
          <code>
            {lines.map((line, index) => (
              <div key={index} className="flex">
                <span className="w-12 shrink-0 text-right pr-4 text-code-line select-none font-mono">
                  {index + 1}
                </span>
                <span className="flex-1 text-foreground font-mono whitespace-pre">
                  {line || " "}
                </span>
              </div>
            ))}
          </code>
        </pre>
      </div>
    </motion.div>
  );
}
