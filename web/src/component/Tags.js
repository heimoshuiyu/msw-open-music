import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

function Tags() {
  const [tags, setTags] = useState([]);
  const [newTagName, setNewTagName] = useState("");
  const [newTagDescription, setNewTagDescription] = useState("");
  const [showAddTag, setShowAddTag] = useState(false);

  function refresh() {
    fetch("/api/v1/get_tags")
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          setTags(data.tags);
        }
      });
  }

  useEffect(() => {
    refresh();
  }, []);

  return (
    <div className="page">
      <h3>Tags</h3>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Description</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {tags.map((tag) => (
            <tr key={tag.id}>
              <td>{tag.name}</td>
              <td>{tag.description}</td>
              <td>
                <Link to={`/manage/tags/${tag.id}`}>Edit</Link>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {!showAddTag && (
        <button onClick={() => setShowAddTag(true)}>Add Tag</button>
      )}
      {showAddTag && (
        <div>
          <label htmlFor="newTagName">New Tag Name</label>
          <input
            type="text"
            id="newTagName"
            value={newTagName}
            onChange={(e) => setNewTagName(e.target.value)}
          />
          <label htmlFor="newTagDescription">New Tag Description</label>
          <textarea
            id="newTagDescription"
            value={newTagDescription}
            onChange={(e) => setNewTagDescription(e.target.value)}
          />
          <button
            onClick={() => {
              fetch("/api/v1/insert_tag", {
                method: "POST",
                headers: {
                  "Content-Type": "application/json",
                },
                body: JSON.stringify({
                  name: newTagName,
                  description: newTagDescription,
                }),
              })
                .then((res) => res.json())
                .then((data) => {
                  if (data.error) {
                    alert(data.error);
                  } else {
                    setNewTagName("");
                    setNewTagDescription("");
                    refresh();
                  }
                });
            }}
          >
            Create Tag
          </button>
        </div>
      )}
    </div>
  );
}

export default Tags;
