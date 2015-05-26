from . import enhancer

class Printer(enhancer.Enhancer):
    def enhance(self, record):
        print(record)
