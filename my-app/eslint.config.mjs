import { dirname } from "path";
import { fileURLToPath } from "url";
import { FlatCompat } from "@eslint/eslintrc";
import nextPlugin from '@eslint/next';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const compat = new FlatCompat({
  baseDirectory: __dirname,
});

const eslintConfig = [
  {
    plugins: {
      next: nextPlugin,
    },
  },
  {
    rules: {
      '@typescript-eslint/no-unused-vars': 'off',
    },
  },
  ...compat.extends("next/core-web-vitals", "next/typescript"),
  ...nextPlugin.configs.recommended,
];

export default eslintConfig;
