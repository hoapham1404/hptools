import { useState, useEffect } from "react";
import { GreetService } from "../bindings/changeme";
import { Events, WML } from "@wailsio/runtime";
import { Button } from "./components/ui/button";

function App() {
  // const [name, setName] = useState<string>('');
  // const [result, setResult] = useState<string>('Please enter your name below 👇');
  // const [time, setTime] = useState<string>('Listening for Time event...');

  // const doGreet = () => {
  //   let localName = name;
  //   if (!localName) {
  //     localName = 'anonymous';
  //   }
  //   GreetService.Greet(localName).then((resultValue: string) => {
  //     setResult(resultValue);
  //   }).catch((err: any) => {
  //     console.log(err);
  //   });
  // }

  // useEffect(() => {
  //   Events.On('time', (timeValue: any) => {
  //     setTime(timeValue.data);
  //   });
  //   // Reload WML so it picks up the wml tags
  //   WML.Reload();
  // }, []);

  return (
    <div className="flex min-h-svh flex-col items-center justify-center">
      <Button>Click me</Button>
    </div>
  );
}

export default App;
