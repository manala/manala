import React, { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from '@app/App';

const container = document.getElementById('app') as HTMLElement;

createRoot(container).render(
  <StrictMode>
    <App />
  </StrictMode>
);
