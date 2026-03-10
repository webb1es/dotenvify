import { Zap, Shield, FileText, GitBranch } from "lucide-react";

const features = [
  { icon: <Zap className="w-3.5 h-3.5" />, text: "Smart quoting for URLs and spaces" },
  { icon: <Shield className="w-3.5 h-3.5" />, text: "0600 perms, auto-backups" },
  { icon: <FileText className="w-3.5 h-3.5" />, text: "Detect missing & unused keys (IDE)" },
  { icon: <GitBranch className="w-3.5 h-3.5" />, text: "Preserve vars across runs" },
];

const FeaturesCell = () => (
  <div className="bento-cell p-4 h-full">
    <p className="font-mono text-[10px] text-muted-foreground mb-3">
      {"// "}what you get
    </p>
    <div className="space-y-3">
      {features.map((f, i) => (
        <div key={i} className="flex items-center gap-2.5">
          <span className="text-primary flex-shrink-0">{f.icon}</span>
          <span className="font-mono text-xs text-foreground">{f.text}</span>
        </div>
      ))}
    </div>
  </div>
);

export default FeaturesCell;
