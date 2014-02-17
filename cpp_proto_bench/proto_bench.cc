#include <cstdlib>
#include <sys/time.h>
#include <iostream>

#include "data.pb.h"

using namespace io::prometheus;

const int numSamples = 5000;
const int numOps = 10000;

uint64_t ClockGetTime()
{
  timespec ts;
  clock_gettime(CLOCK_REALTIME, &ts);
  return (uint64_t)ts.tv_sec * 1000000000LL + (uint64_t)ts.tv_nsec;
}

int main(int argc, char **argv) {
  SampleValueSeries s;

  for (int i = 0; i < numSamples; i++) {
    SampleValueSeries_Value* val = s.add_value();
    val->set_timestamp(rand());
    val->set_value(double(rand()) / double((RAND_MAX)));
  }

  std::string serialized;
  s.SerializeToString(&serialized);

  uint64_t start = ClockGetTime();
  for (int i = 0; i < numOps; i++) {
    SampleValueSeries newSeries;
    newSeries.ParseFromString(serialized);
  }
  uint64_t end = ClockGetTime();
  std::cout << "nanoseconds per op: " << (end - start) / numOps << std::endl;
}
