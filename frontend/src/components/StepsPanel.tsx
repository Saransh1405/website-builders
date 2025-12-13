import { motion } from "framer-motion";
import { Check, Circle, Loader2 } from "lucide-react";

interface Step {
  id: number;
  title: string;
  description: string;
  status: "pending" | "active" | "completed";
}

interface StepsPanelProps {
  steps: Step[];
}

export function StepsPanel({ steps }: StepsPanelProps) {
  return (
    <div className="h-full flex flex-col">
      <div className="p-4 border-b border-border">
        <h2 className="text-lg font-semibold text-foreground">Build Steps</h2>
        <p className="text-sm text-muted-foreground mt-1">
          Progress through each step
        </p>
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-3">
        {steps.map((step, index) => (
          <motion.div
            key={step.id}
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: index * 0.1 }}
            className={`relative p-4 rounded-xl border transition-all ${
              step.status === "active"
                ? "border-primary bg-primary/5"
                : step.status === "completed"
                ? "border-primary/30 bg-primary/5"
                : "border-border bg-card"
            }`}
          >
            <div className="flex items-start gap-3">
              <div
                className={`shrink-0 w-8 h-8 rounded-full flex items-center justify-center ${
                  step.status === "completed"
                    ? "bg-primary text-primary-foreground"
                    : step.status === "active"
                    ? "bg-primary/20 text-primary"
                    : "bg-muted text-muted-foreground"
                }`}
              >
                {step.status === "completed" ? (
                  <Check className="h-4 w-4" />
                ) : step.status === "active" ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <Circle className="h-4 w-4" />
                )}
              </div>

              <div className="flex-1 min-w-0">
                <h3
                  className={`font-medium ${
                    step.status === "pending"
                      ? "text-muted-foreground"
                      : "text-foreground"
                  }`}
                >
                  {step.title}
                </h3>
                <p className="text-sm text-muted-foreground mt-1">
                  {step.description}
                </p>
              </div>
            </div>

            {/* Connection line */}
            {index < steps.length - 1 && (
              <div
                className={`absolute left-[1.9rem] top-14 w-[2px] h-6 ${
                  step.status === "completed" ? "bg-primary/50" : "bg-border"
                }`}
              />
            )}
          </motion.div>
        ))}
      </div>
    </div>
  );
}
