module.exports = {
  env: {
    browser: true,
    es2021: true,
  },
  extends: ["react-app", "react-app/jest", "plugin:react/recommended"],
  parser: "@typescript-eslint/parser",
  parserOptions: {
    ecmaFeatures: {
      jsx: true,
    },
    ecmaVersion: "latest",
    sourceType: "module",
  },
  overrides: [
    {
      files: ["**/*.stories.*"],
      rules: {
        "import/no-anonymous-default-export": "off",
      },
    },
  ],
  plugins: ["react", "@typescript-eslint", "react-hooks"],
  rules: {
    "arrow-body-style": "off",
    "prefer-arrow-callback": "off",
    "react-hooks/rules-of-hooks": "error", // Checks rules of Hooks
    "react-hooks/exhaustive-deps": "warn", // Checks effect dependencies
    // suppress errors for missing 'import React' in files
    "react/react-in-jsx-scope": "off",
    "require-jsdoc": "off",
    "no-unused-vars": "off",
    "@typescript-eslint/no-unused-vars": [
      "warn",
      {
        argsIgnorePattern: "^_",
        varsIgnorePattern: "^_",
        caughtErrorsIgnorePattern: "^_",
      },
    ],
    "no-console": ["warn", { allow: ["warn", "error", "info", "dir"] }],
    // allow jsx syntax in js files (for next.js project)
    "react/jsx-filename-extension": [
      1,
      {
        extensions: [".ts", ".tsx"],
      },
    ],
    "react/prop-types": 0,
  },
  settings: {
    react: {
      version: "detect",
    },
  },
};
