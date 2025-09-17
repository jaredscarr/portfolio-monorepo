export type ColorSchemeMode = "light" | "dark" | "neon" | "ocean" | "forest";

export interface CustomColorScheme {
  mode: ColorSchemeMode;
  colors: {
    primary: string;
    secondary: string;
    background: string;
    surface: string;
    text: {
      primary: string;
      secondary: string;
    };
  };
}

// Example future themes
export const futureThemes: Record<string, CustomColorScheme> = {
  neon: {
    mode: "neon",
    colors: {
      primary: "#ff00ff",
      secondary: "#00ffff",
      background: "#0a0a0a",
      surface: "#1a1a1a",
      text: {
        primary: "#ffffff",
        secondary: "#cccccc",
      },
    },
  },
  ocean: {
    mode: "ocean",
    colors: {
      primary: "#0077be",
      secondary: "#00a86b",
      background: "#f0f8ff",
      surface: "#e6f3ff",
      text: {
        primary: "#003366",
        secondary: "#006699",
      },
    },
  },
  forest: {
    mode: "forest",
    colors: {
      primary: "#228b22",
      secondary: "#8fbc8f",
      background: "#f5f5dc",
      surface: "#f0f8e8",
      text: {
        primary: "#2d4a2d",
        secondary: "#556b2f",
      },
    },
  },
};
