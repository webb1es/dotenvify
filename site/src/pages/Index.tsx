import { Cloud, ClipboardPaste, Search } from "lucide-react";
import Header from "@/components/Header";
import HeroSection from "@/components/HeroSection";
import FeatureCard from "@/components/FeatureCard";
import CodeTransform from "@/components/CodeTransform";
import IdeLogos from "@/components/IdeLogos";
import Footer from "@/components/Footer";

const Index = () => {
  return (
    <>
      <a
        href="#main-content"
        className="sr-only focus:not-sr-only focus:fixed focus:top-2 focus:left-2 focus:z-[100] focus:px-4 focus:py-2 focus:bg-primary focus:text-primary-foreground focus:rounded-lg"
      >
        Skip to main content
      </a>
      <div className="flex flex-col min-h-[100dvh] xl:h-[100dvh] xl:overflow-hidden">
        <Header />

        <main id="main-content" className="flex-1 flex flex-col pt-12">
          {/* Hero */}
          <div className="flex-shrink-0">
            <HeroSection />
          </div>

          {/* Main content grid */}
          <div className="flex-1 min-h-0 px-4 lg:px-6 grid grid-cols-1 lg:grid-cols-[2fr_3fr] gap-3 pb-2">
            {/* Left: Feature cards */}
            <div className="flex flex-col gap-3">
              <FeatureCard
                icon={<Cloud className="w-4 h-4" />}
                title="Azure DevOps Integration"
                description="Fetch variables from Azure DevOps Variable Groups with one click. Secure OAuth via Device Code Flow. Smart merge when .env already exists — per-key conflict resolution."
              />
              <FeatureCard
                icon={<ClipboardPaste className="w-4 h-4" />}
                title="Paste & Format"
                description="Parse any input — KEY=VALUE, quoted, space-separated, line pairs. Smart quoting, alphabetical sorting, export prefix. Live preview as you type."
              />
              <div className="lg:hidden">
                <FeatureCard
                  icon={<Search className="w-4 h-4" />}
                  title="Diagnostics"
                  description="Detect missing and unused env keys across your codebase. Supports JS, Python, Go, Java, Kotlin, Ruby, PHP, Rust, C#, YAML. Click to navigate."
                />
              </div>
            </div>

            {/* Right: Code transform */}
            <div className="min-h-0 hidden md:flex flex-col">
              <CodeTransform />
            </div>
          </div>

          {/* Bottom row */}
          <div className="flex-shrink-0 px-4 lg:px-6 pb-1 grid grid-cols-1 lg:grid-cols-[2fr_3fr] gap-3 items-end">
            <div className="hidden lg:block">
              <FeatureCard
                icon={<Search className="w-4 h-4" />}
                title="Diagnostics"
                description="Detect missing and unused env keys across your codebase. Supports JS, Python, Go, Java, Kotlin, Ruby, PHP, Rust, C#, YAML. Click to navigate."
              />
            </div>
            <IdeLogos />
          </div>
        </main>

        <Footer />
      </div>
    </>
  );
};

export default Index;
