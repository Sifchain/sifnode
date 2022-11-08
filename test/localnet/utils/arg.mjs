import argDep from "arg";

export function arg(
  argObj,
  usageMsg = `This command has no usage message set`
) {
  const args = argDep({ "--help": Boolean, ...argObj });
  if (args["--help"]) {
    console.log(usageMsg);
    process.exit(0);
  }

  return args;
}
