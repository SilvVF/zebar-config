import { createEffect, createSignal, onCleanup, onMount } from "solid-js";
import { createStore } from "solid-js/store";

const RECONNECT_DELAY = 5000;

export const useZbservSocket = (address: string) => {
  const [store, setStore] = createStore({
    connected: false,
    status: {
      genshin: {
        display: "Genshin",
        curr: 0,
        max: 0,
      },
      zzz: {
        display: "ZZZ",
        curr: 0,
        max: 0,
      },
      hkrpg: {
        display: "Star Rail",
        curr: 0,
        max: 0,
      },
    },
  });

  const [conn, setConn] = createSignal<WebSocket | null>(null);

  createEffect(() => {
    const c = conn();
    if (!c) return;

    let reconnectTimeout = undefined;

    const handleMessage = ({ data }: MessageEvent) => {
      try {
        const msg = JSON.parse(data);
        if (msg.details) return;

        const game = msg.game;
        setStore("status", game, { curr: msg.curr, max: msg.max });
      } catch (e) {
        console.error("Failed to parse message", e);
      }
    };

    const handleOpen = () => {
      setStore("connected", true);
    };

    const handleError = () => {
      c.close();
    };

    const handleClose = () => {
      setStore("connected", false);
      reconnectTimeout = setTimeout(() => {
        setConn(() => new WebSocket(address));
      }, RECONNECT_DELAY);
    };

    c.addEventListener("open", handleOpen);
    c.addEventListener("message", handleMessage);
    c.addEventListener("error", handleError);
    c.addEventListener("close", handleClose);

    return () => {
      setStore("connected", false);
      clearTimeout(reconnectTimeout);
      c.removeEventListener("open", handleOpen);
      c.removeEventListener("message", handleMessage);
      c.removeEventListener("error", handleError);
      c.removeEventListener("close", handleClose);
      c.close();
    };
  });

  onMount(async () => {
    // will run only once, when the component is mounted
    setConn(new WebSocket(address));
  });

  onCleanup(() => {
    conn()?.close();
    setConn(null);
  });

  return store;
};
