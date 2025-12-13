import { useState, useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { motion } from "framer-motion";
import { ArrowLeft, Code2 } from "lucide-react";
import { ThemeToggle } from "@/components/ThemeToggle";
import { StepsPanel } from "@/components/StepsPanel";
import { FileExplorer } from "@/components/FileExplorer";
import { CodeViewer } from "@/components/CodeViewer";

interface FileNode {
  name: string;
  type: "file" | "folder";
  children?: FileNode[];
  content?: string;
}

const mockFiles: FileNode[] = [
  {
    name: "src",
    type: "folder",
    children: [
      {
        name: "components",
        type: "folder",
        children: [
          {
            name: "Header.tsx",
            type: "file",
            content: `import React from 'react';

export function Header() {
  return (
    <header className="flex items-center justify-between p-6">
      <div className="text-2xl font-bold">
        Logo
      </div>
      <nav className="flex gap-6">
        <a href="#" className="hover:text-primary">Home</a>
        <a href="#" className="hover:text-primary">About</a>
        <a href="#" className="hover:text-primary">Contact</a>
      </nav>
    </header>
  );
}`,
          },
          {
            name: "Hero.tsx",
            type: "file",
            content: `import React from 'react';

export function Hero() {
  return (
    <section className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <h1 className="text-6xl font-bold mb-6">
          Welcome to Your Site
        </h1>
        <p className="text-xl text-muted-foreground">
          Start building something amazing
        </p>
      </div>
    </section>
  );
}`,
          },
          {
            name: "Footer.tsx",
            type: "file",
            content: `import React from 'react';

export function Footer() {
  return (
    <footer className="border-t border-border py-8">
      <div className="container mx-auto text-center">
        <p className="text-muted-foreground">
          © 2024 Your Company. All rights reserved.
        </p>
      </div>
    </footer>
  );
}`,
          },
        ],
      },
      {
        name: "pages",
        type: "folder",
        children: [
          {
            name: "index.tsx",
            type: "file",
            content: `import { Header } from '@/components/Header';
import { Hero } from '@/components/Hero';
import { Footer } from '@/components/Footer';

export default function Home() {
  return (
    <div className="min-h-screen bg-background">
      <Header />
      <Hero />
      <Footer />
    </div>
  );
}`,
          },
          {
            name: "about.tsx",
            type: "file",
            content: `import { Header } from '@/components/Header';
import { Footer } from '@/components/Footer';

export default function About() {
  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="container mx-auto py-12">
        <h1 className="text-4xl font-bold mb-6">About Us</h1>
        <p className="text-muted-foreground">
          Learn more about our company and mission.
        </p>
      </main>
      <Footer />
    </div>
  );
}`,
          },
        ],
      },
      {
        name: "styles",
        type: "folder",
        children: [
          {
            name: "globals.css",
            type: "file",
            content: `@tailwind base;
@tailwind components;
@tailwind utilities;

:root {
  --background: 0 0% 100%;
  --foreground: 222.2 84% 4.9%;
  --primary: 168 84% 40%;
}

.dark {
  --background: 220 20% 6%;
  --foreground: 220 10% 95%;
  --primary: 168 84% 50%;
}

body {
  font-family: system-ui, sans-serif;
}`,
          },
        ],
      },
    ],
  },
  {
    name: "public",
    type: "folder",
    children: [
      {
        name: "favicon.ico",
        type: "file",
        content: "Binary file - favicon icon",
      },
      {
        name: "robots.txt",
        type: "file",
        content: `User-agent: *
Allow: /

Sitemap: https://example.com/sitemap.xml`,
      },
    ],
  },
  {
    name: "package.json",
    type: "file",
    content: `{
  "name": "my-website",
  "version": "1.0.0",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "^18.3.1",
    "react-dom": "^18.3.1"
  },
  "devDependencies": {
    "vite": "^5.0.0",
    "tailwindcss": "^3.4.0"
  }
}`,
  },
  {
    name: "tailwind.config.js",
    type: "file",
    content: `/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class',
  content: ['./src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        primary: 'hsl(var(--primary))',
      },
    },
  },
  plugins: [],
}`,
  },
];

const initialSteps = [
  {
    id: 1,
    title: "Analyzing prompt",
    description: "Understanding your requirements and planning the structure",
    status: "completed" as const,
  },
  {
    id: 2,
    title: "Creating project structure",
    description: "Setting up folders, files, and configurations",
    status: "completed" as const,
  },
  {
    id: 3,
    title: "Building components",
    description: "Generating React components for your website",
    status: "active" as const,
  },
  {
    id: 4,
    title: "Adding styles",
    description: "Applying beautiful, responsive design",
    status: "pending" as const,
  },
  {
    id: 5,
    title: "Final touches",
    description: "Optimizing and polishing the final result",
    status: "pending" as const,
  },
];

const Builder = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const prompt = location.state?.prompt || "Create a modern website";
  const [selectedFile, setSelectedFile] = useState<FileNode | null>(null);
  const [steps, setSteps] = useState(initialSteps);

  useEffect(() => {
    // Simulate step progression
    const timers: NodeJS.Timeout[] = [];
    
    timers.push(
      setTimeout(() => {
        setSteps((prev) =>
          prev.map((step) =>
            step.id === 3
              ? { ...step, status: "completed" as const }
              : step.id === 4
              ? { ...step, status: "active" as const }
              : step
          )
        );
      }, 3000)
    );

    timers.push(
      setTimeout(() => {
        setSteps((prev) =>
          prev.map((step) =>
            step.id === 4
              ? { ...step, status: "completed" as const }
              : step.id === 5
              ? { ...step, status: "active" as const }
              : step
          )
        );
      }, 6000)
    );

    timers.push(
      setTimeout(() => {
        setSteps((prev) =>
          prev.map((step) =>
            step.id === 5 ? { ...step, status: "completed" as const } : step
          )
        );
      }, 9000)
    );

    return () => timers.forEach(clearTimeout);
  }, []);

  return (
    <div className="h-screen flex flex-col bg-background">
      {/* Header */}
      <header className="flex items-center justify-between px-4 py-3 border-b border-border bg-card">
        <div className="flex items-center gap-4">
          <button
            onClick={() => navigate("/")}
            className="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Back
          </button>
          <div className="h-6 w-px bg-border" />
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
              <Code2 className="h-4 w-4 text-primary-foreground" />
            </div>
            <span className="font-medium text-foreground">BuilderAI</span>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <div className="px-3 py-1.5 rounded-lg bg-secondary text-sm text-secondary-foreground max-w-md truncate">
            {prompt}
          </div>
          <ThemeToggle />
        </div>
      </header>

      {/* Main content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Steps panel */}
        <motion.aside
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          className="w-80 border-r border-border bg-card shrink-0 overflow-hidden"
        >
          <StepsPanel steps={steps} />
        </motion.aside>

        {/* Code viewer */}
        <main className="flex-1 overflow-hidden bg-background">
          {selectedFile ? (
            <CodeViewer
              fileName={selectedFile.name}
              content={selectedFile.content || ""}
            />
          ) : (
            <div className="h-full flex items-center justify-center text-muted-foreground">
              <p>Select a file to view its contents</p>
            </div>
          )}
        </main>

        {/* File explorer */}
        <motion.aside
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          className="w-72 border-l border-border bg-card shrink-0 overflow-hidden"
        >
          <FileExplorer
            files={mockFiles}
            onFileSelect={setSelectedFile}
            selectedFile={selectedFile}
          />
        </motion.aside>
      </div>
    </div>
  );
};

export default Builder;
