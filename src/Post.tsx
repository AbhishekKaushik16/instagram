import React from "react";
import "./styles/Post.css";
import reactIcon from "./assets/react.svg";
import { Avatar } from "@mui/material";

const Post = () => {
  return (
    <div className="post">
      <Avatar
        className="post__avatar"
        alt={"Abhishek"}
        src="/static/images/avatar/1.jpg"
      />
      <h3>Username</h3>
      <img src={reactIcon} className="post__image" alt="React logo" />
    </div>
  );
};

export default Post;
