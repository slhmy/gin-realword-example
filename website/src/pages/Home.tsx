import { useState } from "react";
import { useChatStream } from "@/hooks/chatStream";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import MarkdownRender from "@/components/MarkdownRender";
import { useUser } from "@/hooks/user";

function Home() {
  const [prompt, setPrompt] = useState("");
  const { user } = useUser();
  const { createChatStream, isLoading, response } = useChatStream();

  return (
    <div className="flex flex-col gap-4 max-w-2xl mx-auto p-16">
      <h1 className="text-2xl font-bold mb-4">Chat with AI</h1>
      <Textarea onChange={(e) => setPrompt(e.target.value)} />
      {!user && (
        <div className="text-red-500">
          Please{" "}
          <a href="/auth/github/login" className="text-blue-500 underline">
            login
          </a>{" "}
          to use this feature.
        </div>
      )}
      <Button
        className="mb-8"
        onClick={() => createChatStream(prompt)}
        disabled={isLoading || !user}
      >
        {isLoading ? "Loading..." : "Submit"}
      </Button>
      <MarkdownRender content={response} />
    </div>
  );
}

export default Home;
