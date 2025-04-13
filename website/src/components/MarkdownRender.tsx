import React from "react";
import ReactMarkdown from "react-markdown";

interface MarkdownRenderProps {
  content: string | null;
}

const MarkdownRender: React.FC<MarkdownRenderProps> = (props) => {
  return (
    <article className="prose">
      <ReactMarkdown>{props.content}</ReactMarkdown>
    </article>
  );
};

export default MarkdownRender;
