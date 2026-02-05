import { useEffect, useMemo, useState } from "react";
import { WindowMinimise, WindowToggleMaximise, Quit } from "../wailsjs/runtime";

declare global {
  interface Window {
    go: any;
  }
}

type UIState = {
  IsActive: boolean;
  ZoomLevel: number;
  ScopeSize: number;
  BorderEnabled: boolean;
  BorderColor: number;
  BorderThickness: number;
  FollowCursor: boolean;
  HotkeyMode: string;
  Hotkey: number;
  RenderBackend: string;
};

function clamp(min: number, max: number, v: number) {
  return Math.min(max, Math.max(min, v));
}

export default function App() {
  const [state, setState] = useState<UIState>({
    IsActive: false,
    ZoomLevel: 2.0,
    ScopeSize: 250,
    BorderEnabled: true,
    BorderColor: 0,
    BorderThickness: 2,
    FollowCursor: false,
    HotkeyMode: "Toggle",
    Hotkey: 0,
    RenderBackend: "Auto",
  });

  const [rainbowEnabled, setRainbowEnabled] = useState(false);
  const [rainbowSpeed, setRainbowSpeed] = useState(1.5);

  const hotkeyOptions = [
    { name: "Nenhum", code: 0x00 },
    { name: "Clique Direito (Mouse)", code: 0x02 },
    { name: "Botão Lateral 1 (Mouse)", code: 0x05 },
    { name: "Botão Lateral 2 (Mouse)", code: 0x06 },
    { name: "Botão do Meio (Mouse)", code: 0x04 },
    { name: "Shift Esquerdo", code: 0xA0 },
    { name: "Ctrl Esquerdo", code: 0xA2 },
    { name: "Alt Esquerdo", code: 0xA4 },
    { name: "Caps Lock", code: 0x14 },
    { name: "V", code: 0x56 },
    { name: "C", code: 0x43 },
    { name: "F", code: 0x46 },
    { name: "X", code: 0x58 },
    { name: "Z", code: 0x5A },
  ];

  useEffect(() => {
    async function load() {
      try {
        const s = await window.go.wailsapp.App.GetState();
        setState(s);
      } catch {
        // noop: during dev without wails runtime
      }
    }
    load();
  }, []);

  const setActive = async (on: boolean) => {
    setState((s) => ({ ...s, IsActive: on }));
    await window.go.wailsapp.App.SetActive(on);
  };

  const toggleActive = async () => {
    const next = !state.IsActive;
    setState((s) => ({ ...s, IsActive: next }));
    await window.go.wailsapp.App.ToggleActive();
  };

  const setHotkey = async (code: number) => {
    setState((s) => ({ ...s, Hotkey: code }));
    await window.go.wailsapp.App.SetHotkey(code);
  };

  const setHotkeyMode = async (mode: string) => {
    const val = mode === "Hold" ? "Hold" : "Toggle";
    setState((s) => ({ ...s, HotkeyMode: val }));
    await window.go.wailsapp.App.SetHotkeyMode(val);
  };

  const setZoom = async (v: number) => {
    const val = Math.round(clamp(1.0, 8.0, v) * 10) / 10;
    setState((s) => ({ ...s, ZoomLevel: val }));
    await window.go.wailsapp.App.SetZoomLevel(val);
  };

  const setBorderEnabled = async (v: boolean) => {
    setState((s) => ({ ...s, BorderEnabled: v }));
    await window.go.wailsapp.App.SetBorderEnabled(v);
    if (!v) {
      setRainbowEnabled(false);
    }
  };

  const setBorderThickness = async (v: number) => {
    const val = Math.round(clamp(1, 10, v));
    setState((s) => ({ ...s, BorderThickness: val }));
    await window.go.wailsapp.App.SetBorderThickness(val);
  };

  const pickColorIndex = (hex: string) => {
    const h = hex.replace("#", "").toLowerCase();
    const colors = ["ff0000", "00ff00", "0000ff", "ffff00"];
    let idx = 0;
    let best = 1e9;
    const r = parseInt(h.slice(0, 2), 16);
    const g = parseInt(h.slice(2, 4), 16);
    const b = parseInt(h.slice(4, 6), 16);
    for (let i = 0; i < colors.length; i++) {
      const cr = parseInt(colors[i].slice(0, 2), 16);
      const cg = parseInt(colors[i].slice(2, 4), 16);
      const cb = parseInt(colors[i].slice(4, 6), 16);
      const dist = (r - cr) ** 2 + (g - cg) ** 2 + (b - cb) ** 2;
      if (dist < best) {
        best = dist;
        idx = i;
      }
    }
    return idx;
  };

  const setBorderColorByHex = async (hex: string) => {
    const idx = pickColorIndex(hex);
    setState((s) => ({ ...s, BorderColor: idx }));
    await window.go.wailsapp.App.SetBorderColor(idx);
    setRainbowEnabled(false);
  };

  const setScopeSize = async (px: number) => {
    const val = Math.round(clamp(100, 600, px));
    setState((s) => ({ ...s, ScopeSize: val }));
    await window.go.wailsapp.App.SetScopeSize(val);
  };

  const setFollowCursor = async (on: boolean) => {
    setState((s) => ({ ...s, FollowCursor: on }));
    await window.go.wailsapp.App.SetFollowCursor(on);
  };

  const setRenderBackend = async (label: string) => {
    const value = label === "Compatibilidade (GDI)" ? "GDI" : "Auto";
    setState((s) => ({ ...s, RenderBackend: value }));
    await window.go.wailsapp.App.SetRenderBackend(value);
  };

  // Rainbow effect: changes border color index periodically
  useEffect(() => {
    if (!rainbowEnabled || !state.BorderEnabled) return;
    let mounted = true;
    const interval = Math.max(0.2, rainbowSpeed); // changes per second
    const ms = Math.round(1000 / interval);
    const id = setInterval(async () => {
      if (!mounted) return;
      const next = Math.floor(Math.random() * 4);
      setState((s) => ({ ...s, BorderColor: next }));
      await window.go.wailsapp.App.SetBorderColor(next);
    }, ms);
    return () => {
      mounted = false;
      clearInterval(id);
    };
  }, [rainbowEnabled, rainbowSpeed, state.BorderEnabled]);

  const rainbowDisabled = useMemo(() => !state.BorderEnabled, [state.BorderEnabled]);

  return (
    <div className="p-2 space-y-4">
      <div className="titlebar wails-draggable">
        <div className="text-sm font-medium text-gray-300">Mira Controller</div>
        <div className="buttons wails-nodrag">
          <button className="btn-icon" title="Minimizar" onClick={() => WindowMinimise()}>
            ─
          </button>
          <button className="btn-icon" title="Maximizar" onClick={() => WindowToggleMaximise()}>
            ☐
          </button>
          <button className="btn-icon" title="Fechar" onClick={() => Quit()}>
            ✕
          </button>
        </div>
      </div>
      <div className="grid grid-cols-2 gap-4">
      <div className="panel p-4 space-y-3">
        <div className="flex items-center justify-between">
          <div className="title">Mira Controller</div>
          <button
            className={`px-3 py-1 rounded-md text-sm font-medium ${
              state.IsActive
                ? "bg-primary text-white"
                : "bg-gray-700 text-gray-200"
            }`}
            onClick={toggleActive}
          >
            {state.IsActive ? "Luneta: Ligada" : "Luneta: Desligada"}
          </button>
        </div>
        <p className="label">Painel de controle moderno para a lunetinha</p>
      </div>

      <div className="panel p-4 space-y-4">
        <div className="title mb-2">Ativação</div>
        <div className="grid grid-cols-2 gap-4">
          <div className="flex flex-col">
            <span className="label">Tecla/Botão</span>
            <select
              className="mt-2 rounded-md bg-gray-800 border border-gray-700 px-2 py-2"
              value={state.Hotkey}
              onChange={(e) => setHotkey(parseInt(e.target.value))}
            >
              {hotkeyOptions.map((k) => (
                <option key={k.code} value={k.code}>
                  {k.name}
                </option>
              ))}
            </select>
          </div>
          <div className="flex flex-col">
            <span className="label">Modo</span>
            <select
              className="mt-2 rounded-md bg-gray-800 border border-gray-700 px-2 py-2"
              value={state.HotkeyMode}
              onChange={(e) => setHotkeyMode(e.target.value)}
            >
              <option value="Toggle">Alternar (Toggle)</option>
              <option value="Hold" disabled={state.Hotkey === 0}>Segurar (Hold)</option>
            </select>
          </div>
        </div>
      </div>

      <div className="panel p-4 space-y-4">
        <div>
          <div className="flex justify-between">
            <span className="label">Nível de zoom</span>
            <span className="label">{state.ZoomLevel.toFixed(1)}x</span>
          </div>
          <input
            className="slider mt-2"
            type="range"
            min={10}
            max={80}
            value={Math.round(state.ZoomLevel * 10)}
            onChange={(e) => setZoom(parseInt(e.target.value) / 10)}
          />
        </div>

        <div>
          <div className="flex justify-between">
            <span className="label">Tamanho da luneta</span>
            <span className="label">{Math.round(state.ScopeSize)} px</span>
          </div>
          <input
            className="slider mt-2"
            type="range"
            min={100}
            max={600}
            value={Math.round(state.ScopeSize)}
            onChange={(e) => setScopeSize(parseInt(e.target.value))}
          />
        </div>
        <div className="flex items-center justify-between">
          <span className="label">Seguir cursor</span>
          <label className="inline-flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={state.FollowCursor}
              onChange={(e) => setFollowCursor(e.target.checked)}
            />
            <span className="label">Ativar</span>
          </label>
        </div>
        <div className="flex items-center justify-between">
          <span className="label">Backend de renderização</span>
          <select
            className="mt-2 rounded-md bg-gray-800 border border-gray-700 px-2 py-2"
            value={
              state.RenderBackend === "GDI"
                ? "Compatibilidade (GDI)"
                : "Automático (DirectX)"
            }
            onChange={(e) => setRenderBackend(e.target.value)}
          >
            <option>Automático (DirectX)</option>
            <option>Compatibilidade (GDI)</option>
          </select>
        </div>
      </div>

      <div className="panel p-4 space-y-4">
        <div className="flex items-center justify-between">
          <span className="title">Borda</span>
          <label className="inline-flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={state.BorderEnabled}
              onChange={(e) => setBorderEnabled(e.target.checked)}
            />
            <span className="label">Ativar</span>
          </label>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <div className="flex justify-between">
              <span className="label">Espessura</span>
              <span className="label">{Math.round(state.BorderThickness)}</span>
            </div>
            <input
              className="slider mt-2"
              type="range"
              min={1}
              max={10}
              value={Math.round(state.BorderThickness)}
              onChange={(e) => setBorderThickness(parseInt(e.target.value))}
              disabled={!state.BorderEnabled}
            />
          </div>

          <div className="flex flex-col">
            <span className="label">Cor da borda</span>
            <input
              type="color"
              className="mt-2 h-10 w-16 p-1 rounded-md bg-gray-800 border border-gray-700"
              onChange={(e) => setBorderColorByHex(e.target.value)}
              value={
                ["#ff0000", "#00ff00", "#0000ff", "#ffff00"][state.BorderColor]
              }
              disabled={!state.BorderEnabled || rainbowEnabled}
            />
          </div>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div className="flex items-center gap-3">
            <label className="inline-flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                checked={rainbowEnabled}
                onChange={(e) => setRainbowEnabled(e.target.checked)}
                disabled={rainbowDisabled}
              />
              <span className="label">Borda Rainbow</span>
            </label>
          </div>
          <div>
            <div className="flex justify-between">
              <span className="label">Velocidade</span>
              <span className="label">{rainbowSpeed.toFixed(1)}x/s</span>
            </div>
            <input
              className="slider mt-2"
              type="range"
              min={2}
              max={20}
              value={Math.round(rainbowSpeed * 10)}
              onChange={(e) => setRainbowSpeed(parseInt(e.target.value) / 10)}
              disabled={!rainbowEnabled}
            />
          </div>
        </div>
      </div>

      {/* <div className="panel p-4">
        <div className="flex items-center justify-between">
          <span className="label">Ligar/Desligar</span>
          <button
            className="px-3 py-2 rounded-md bg-primary hover:opacity-90 transition"
            onClick={() => setActive(!state.IsActive)}
          >
            {state.IsActive ? "Desligar" : "Ligar"}
          </button>
        </div>
      </div> */}
      </div>
    </div>
  );
}
