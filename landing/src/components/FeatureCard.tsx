import { cn } from "@/lib/utils";

interface FeatureCardProps {
  icon: React.ReactNode;
  title: string;
  description: string;
  className?: string;
}

const FeatureCard = ({ icon, title, description, className }: FeatureCardProps) => (
  <article
    className={cn(
      "rounded-lg border border-jb-surface-border bg-card p-4 card-hover-glow",
      className
    )}
  >
    <div className="flex items-start gap-3">
      <div className="flex-shrink-0 w-8 h-8 rounded-md bg-jb-blue/10 flex items-center justify-center text-jb-blue">
        {icon}
      </div>
      <div className="min-w-0">
        <h2 className="text-sm font-semibold text-foreground">{title}</h2>
        <p className="text-xs text-muted-foreground mt-1 leading-relaxed">{description}</p>
      </div>
    </div>
  </article>
);

export default FeatureCard;
