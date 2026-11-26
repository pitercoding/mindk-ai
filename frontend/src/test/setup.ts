// Adds the jest-dom matchers (toBeInTheDocument, toHaveValue, ...).
// React Testing Library unmounts components after each test by itself because
// vitest `globals` is enabled, so no manual cleanup is needed here.
import "@testing-library/jest-dom/vitest";

// jsdom does not implement scrolling; NotesPage calls window.scrollTo.
window.scrollTo = () => {};
