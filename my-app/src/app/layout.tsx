import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Glitchy Email Form",
  description: "A minimal and glitchy email form with neon styling",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}
