import { Link } from "react-router-dom";
import { convertIntToDateTime } from "./Common";

function ReviewEntry(props) {
  return (
    <div>
      <h4>
        <Link to={`/manage/users/${props.review.user.id}`}>
          @{props.review.user.username}
        </Link>{" "}
        review{" "}
        <Link to={`/files/${props.review.file.id}`}>
          {props.review.file.filename}
        </Link>{" "}
        on {convertIntToDateTime(props.review.created_at)}{" "}
        {props.review.updated_at !== 0 &&
          "(modified on " +
            convertIntToDateTime(props.review.updated_at) +
            ")"}{" "}
        {(props.user.role === 1 || props.review.user.id === props.user.id) &&
          props.user.role !== 0 && (
            <Link to={`/manage/reviews/${props.review.id}`}>Edit</Link>
          )}
      </h4>
      <p>{props.review.content}</p>
    </div>
  );
}

export default ReviewEntry;
