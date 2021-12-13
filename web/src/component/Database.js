import { useState, useEffect } from "react";

function Database() {
  const [walkPath, setWalkPath] = useState("");
  const [patternString, setPatternString] = useState("");
  const [tags, setTags] = useState([]);
  const [selectedTags, setSelectedTags] = useState([]);

  function getTags() {
    fetch("/api/v1/get_tags")
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          setTags(data.tags);
        }
      });
  }

  useEffect(() => {
    getTags();
  }, []);

  function updateDatabase() {
    // split pattern string into array
    let patternArray = patternString.split(" ");
    // remove whitespace from array
    patternArray = patternArray.map((item) => item.trim());
    // remove empty strings from array
    patternArray = patternArray.filter((item) => item !== "");
    // add dot before item array
    patternArray = patternArray.map((item) => "." + item);

    fetch("/api/v1/walk", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        root: walkPath,
        pattern: patternArray,
        tag_ids: selectedTags,
      }),
    })
      .then((res) => res.json())
      .then((data) => {
        console.log(data);
      });
  }
  return (
    <div>
      <h3>Update Database</h3>
      <input
        type="text"
        value={walkPath}
        placeholder="walk path"
        onChange={(e) => setWalkPath(e.target.value)}
      />
      <input
        type="text"
        value={patternString}
        placeholder="pattern wav flac mp3"
        onChange={(e) => setPatternString(e.target.value)}
      />
      <div>
        <h4>Tags</h4>
        {tags.map((tag) => (
          <div key={tag.id}>
            <input
              id={tag.id}
              type="checkbox"
              value={tag.id}
              onChange={(e) => {
                if (e.target.checked) {
                  setSelectedTags([...selectedTags, tag.id]);
                } else {
                  setSelectedTags(
                    selectedTags.filter((item) => item !== tag.id)
                  );
                }
              }}
            />
            <label htmlFor={tag.id}>{tag.name}</label>
          </div>
        ))}
      </div>
      <button
        onClick={() => {
          updateDatabase();
        }}
      >
        Update Database
      </button>
    </div>
  );
}

export default Database;
