/* @refresh reload */
import "./index.css";
import { For, render } from "solid-js/web";
import { createStore } from "solid-js/store";
import * as zebar from "zebar";
import { useZbservSocket } from "./hoyoStore";

const providers = zebar.createProviderGroup({
  audio: { type: "audio" },
  cpu: { type: "cpu" },
  battery: { type: "battery" },
  memory: { type: "memory" },
  weather: { type: "weather" },
  media: { type: "media" },
  systray: { type: "systray" },
});

render(() => <App />, document.getElementById("root")!);

function App() {
  const [output, setOutput] = createStore(providers.outputMap);

  providers.onOutput((outputMap) => setOutput(outputMap));

  const zbstore = useZbservSocket("ws://localhost:45456/ws");

  return (
    <div class="flex flex-row h-full w-full items-center justify-end overflow-clip text-center text-foreground">
      {zbstore.connected && (
        <For each={Object.values(zbstore.status)}>
          {(s) => (
            <div class="inline-block px-2 py-1 rounded bg-background mr-1">
              {s.display}: {s.curr}/{s.max}
            </div>
          )}
        </For>
      )}
      {output.media?.currentSession && (
        <div class="inline-block px-2 py-1 rounded bg-background mr-1">
          {output.media.currentSession.title}-
          {output.media.currentSession.artist}
          <button onClick={() => output.media?.togglePlayPause()}>⏯</button>
        </div>
      )}
      {output.audio?.defaultPlaybackDevice && (
        <div class="flex flex-row items-center px-2 py-1 rounded bg-background mr-1 space-x-2">
          <div>%{output.audio.defaultPlaybackDevice.volume}</div>
          <input
            type="range"
            min="0"
            max="100"
            step="2"
            value={output.audio.defaultPlaybackDevice.volume}
            onChange={(e) => output.audio.setVolume(e.target.valueAsNumber)}
          />
        </div>
      )}
      {output.cpu && (
        <div class="inline-block px-2 py-1 rounded bg-background mr-1">
          CPU: {Math.round(output.cpu.usage)}
        </div>
      )}
      {output.memory && (
        <div class="inline-block px-2 py-1 rounded bg-background mr-1">
          Memory: {Math.round(output.memory.usage)}
        </div>
      )}
      {output.weather && (
        <div class="inline-block px-2 py-1 rounded bg-background mr-1">
          Temp: {Math.round(output.weather.celsiusTemp)}
        </div>
      )}
      {output.systray && (
        <div class="flex flex-row px-2 py-1 rounded bg-background mr-1">
          <For each={output.systray.icons}>
            {(icon) => (
              <img
                class="w-4 h-4 mr-1"
                src={icon.iconUrl}
                title={icon.tooltip}
                onClick={(e) => {
                  e.preventDefault();
                  output.systray.onLeftClick(icon.id);
                }}
                onContextMenu={(e) => {
                  e.preventDefault();
                  output.systray.onRightClick(icon.id);
                }}
              />
            )}
          </For>
        </div>
      )}
    </div>
  );
}
