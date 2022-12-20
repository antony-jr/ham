import React from 'react';

export default function DemoVideo({video}) {
  return (
     <div style={{width: "100%"}}>
	<center>
	<video width="90%" height="" controls={true} muted={true} autoPlay={true} loop={true}>
	   <source src={require("@site/static/vids/" + video + ".webm").default} type="video/webm" />
	</video>
	</center>
     </div>
  );
}
