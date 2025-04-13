import { useState } from "react";

const MAX_TIMEOUT = 300000; // 5 minutes

export const useChatStream = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [response, setResponse] = useState<string | null>(null);

  function createChatStream(prompt: string) {
    setResponse(null);
    setIsLoading(true);
    const encodedPrompt = encodeURIComponent(prompt);
    const eventSource = new EventSource(
      `/api/v1/chat/stream?prompt=${encodedPrompt}`
    );
    const closeEventSource = () => {
      eventSource.close();
      setIsLoading(false);
    };
    eventSource.onopen = function (event) {
      console.log(event);
    };
    eventSource.onmessage = function (event) {
      console.log(event);
      const data = event.data;
      const parsedData = JSON.parse(data);
      if (parsedData.done) {
        closeEventSource();
        return;
      }
      const message = parsedData.content;
      setResponse((prev) => (prev ? prev + message : message));
    };
    eventSource.onerror = function (event) {
      console.error(event);
      closeEventSource();
    };

    setTimeout(() => {
      if (eventSource.readyState !== EventSource.CLOSED) {
        closeEventSource();
      }
    }, MAX_TIMEOUT);
  }

  return { createChatStream, isLoading, response };
};
