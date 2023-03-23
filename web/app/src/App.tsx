import React, { StrictMode } from "react";
import { createRoot } from "react-dom/client";

const container = document.getElementById('app');

function App() {
    return <div>
        Manala Web UI
    </div>
}

createRoot(container).render(
  <StrictMode>
    <App />
  </StrictMode>
);
