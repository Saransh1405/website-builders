import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Sparkles, ArrowRight } from "lucide-react";
import { motion } from "framer-motion";

export function PromptInput() {
  const [prompt, setPrompt] = useState("");
  const [isFocused, setIsFocused] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (prompt.trim()) {
      navigate("/builder", { state: { prompt: prompt.trim() } });
    }
  };

  return (
    <motion.form
      onSubmit={handleSubmit}
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.3, duration: 0.5 }}
      className="w-full max-w-3xl mx-auto"
    >
      <div
        className={`relative rounded-2xl transition-all duration-300 ${
          isFocused ? "glow-primary" : ""
        }`}
      >
        {/* Gradient border */}
        <div className="absolute inset-0 rounded-2xl bg-gradient-to-r from-primary via-primary/50 to-primary opacity-100 p-[1px]">
          <div className="w-full h-full rounded-2xl bg-card" />
        </div>

        {/* Input container */}
        <div className="relative flex items-center gap-3 p-2 pl-5">
          <Sparkles className="h-5 w-5 text-primary shrink-0" />
          
          <input
            type="text"
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            onFocus={() => setIsFocused(true)}
            onBlur={() => setIsFocused(false)}
            placeholder="Describe the website you want to create..."
            className="flex-1 bg-transparent text-foreground placeholder:text-muted-foreground focus:outline-none text-lg py-4"
          />

          <button
            type="submit"
            disabled={!prompt.trim()}
            className="shrink-0 flex items-center gap-2 px-6 py-3 rounded-xl bg-primary text-primary-foreground font-medium transition-all hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Create
            <ArrowRight className="h-4 w-4" />
          </button>
        </div>
      </div>

      <p className="text-center text-sm text-muted-foreground mt-4">
        Press Enter or click Create to start building
      </p>
    </motion.form>
  );
}
