import { Terminal, Puzzle, Code2 } from "lucide-react";
import Header from "@/components/Header";
import HeroCell from "@/components/HeroCell";
import LiveDemo from "@/components/LiveDemo";
import ProductCard from "@/components/ProductCard";
import FormatsCell from "@/components/FormatsCell";
import FeaturesCell from "@/components/FeaturesCell";
import Footer from "@/components/Footer";

const Index = () => (
  <>
    <a
      href="#main-content"
      className="sr-only focus:not-sr-only focus:fixed focus:top-2 focus:left-2 focus:z-[100] focus:px-4 focus:py-2 focus:bg-primary focus:text-primary-foreground focus:rounded-lg"
    >
      Skip to main content
    </a>

    <div className="min-h-[100dvh] flex flex-col">
      <Header />

      <main id="main-content" className="flex-1 px-3 lg:px-4 pb-3 lg:pb-4">
        <div className="max-w-6xl mx-auto grid gap-3 grid-cols-1 lg:grid-cols-12 auto-rows-auto">

          {/* Row 1: Hero + Demo */}
          <div className="lg:col-span-4 lg:row-span-2 min-h-[280px]">
            <HeroCell />
          </div>
          <div className="lg:col-span-8 lg:row-span-2 min-h-[280px] lg:min-h-[400px]">
            <LiveDemo />
          </div>

          {/* Row 2: Products */}
          <div className="lg:col-span-4">
            <ProductCard
              name="CLI"
              icon={<Terminal className="w-4 h-4" />}
              description="Pipe-friendly. Scriptable. Convert any messy input into a clean .env from the terminal."
              installCmd="npx dotenvify input.txt -o .env"
              href="https://github.com/webb1es/dotenvify"
              ctaLabel="view source"
            />
          </div>
          <div className="lg:col-span-4">
            <ProductCard
              name="JetBrains"
              icon={<Puzzle className="w-4 h-4" />}
              description="Pull Azure DevOps variables, paste & format, diagnostics — without leaving your IDE."
              href="https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify"
              ctaLabel="get plugin"
              badges={["IntelliJ", "GoLand", "WebStorm", "PyCharm", "Rider"]}
            />
          </div>
          <div className="lg:col-span-4">
            <ProductCard
              name="VS Code"
              icon={<Code2 className="w-4 h-4" />}
              description="Same powerful features, coming to Visual Studio Code. Parser, formatter, diagnostics — all built on @dotenvify/core."
              comingSoon
            />
          </div>

          {/* Row 3: Formats + Features */}
          <div className="lg:col-span-5">
            <FormatsCell />
          </div>
          <div className="lg:col-span-7">
            <FeaturesCell />
          </div>
        </div>
      </main>

      <Footer />
    </div>
  </>
);

export default Index;
